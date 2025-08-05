import http from "k6/http";
import exec from "k6/execution";
import { check } from "k6";
import {
  randomString,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { SharedArray } from "k6/data";

const domain = __ENV.DOMAIN || "localhost:8080"
const number_visitors = __ENV.N_VISITORS || 10
const number_pages = __ENV.N_PAGES || 10
const number_visits_per_visitor = __ENV.N_VISITS || 10

export function setup() {
  console.log("Domain:                       ", domain)
  console.log("Number of Visitors:           ", number_visitors)
  console.log("Number of Pages:              ", number_pages)
  console.log("Number of Visits per Visitor: ", number_visits_per_visitor)
}

// generates unique visitor ids
const visitors = new SharedArray("visitors", function () {
  var arr = [];
  for (var i = 0; i < number_visitors; i++) {
    arr.push("visitor-" + i);
  }
  return arr;
});

// generates unique page urls
const pages = new SharedArray("pages", function () {
  var arr = [];
  for (var i = 0; i < number_pages; i++) {
    arr.push(randomString(8));
  }
  return arr;
});

// generates random data based on the pages and visitors available
const requests = new SharedArray("requests", function () {
  var arr = [];
  for (var i = 0; i < number_visitors * number_visits_per_visitor; i++) {
    arr.push({
      visitor_id: randomItem(visitors),
      page_url: randomItem(pages),
    });
  }
  return arr;
});

// processes the data to know what to expect when asking the server for unique number of visitors to a page
const expected = new SharedArray("expected", function () {
  var map = new Map();
  for (var i = 0; i < pages.length; i++) {
    map.set(pages[i], new Set());
  }
  for (var i = 0; i < requests.length; i++) {
    var s = map.get(requests[i].page_url);
    s.add(requests[i].visitor_id);
    map.set(requests[i].page_url, s);
  }
  var arr = [];
  for (const [key, value] of map) {
    arr.push({
      counter: value.size,
      page_url: key,
    });
  }
  return arr;
});

export const options = {
  //number of visitors defines how many concurrent clients will call the service
  vus: number_visitors,
  iterations: number_visitors * number_visits_per_visitor,
};

export default function () {
  const url = `http://${domain}/api/v1/user-navigation`;
  const payload = JSON.stringify(requests[exec.scenario.iterationInTest]);

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  http.post(url, payload, params);
}

// verifies if the server calculated the data correctly and is safe from data races
export function teardown(data) {
  for (var i = 0; i < expected.length; i++) {
    const url =
      `http://${domain}/api/v1/unique-visitors?pageUrl=` +
      expected[i].page_url;
    const res = http.get(url);
    check(res, {
      "is counter correct": (r) =>
        r.body.includes(`"unique_visitors":` + expected[i].counter),
    });
  }
}
