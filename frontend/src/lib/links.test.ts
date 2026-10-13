import test from "node:test";
import assert from "node:assert/strict";
import { filterLinks, sortLinks, readView, requestJSON, type Link } from "./links.ts";

const links: Link[] = [
 {short:"zebra",long:"https://example.org/one",created:"",lastEdit:"",clicks:3},
 {short:"Alpha",long:"https://wiki.test/",created:"",lastEdit:"",clicks:9},
];
test("search uses names and destinations, sort keeps input intact",()=>{
 assert.deepEqual(filterLinks(links,"WIKI"),[links[1]]);
 assert.equal(sortLinks(links,"name")[0].short,"Alpha");
 assert.equal(sortLinks(links,"clicks")[0].short,"Alpha");
 assert.equal(links[0].short,"zebra");
});
test("URL preserves search/editor state and unknown link creation",()=>{
 assert.equal(readView(new URL("http://go/?edit=Alpha&q=wiki")).edit,"Alpha");
 assert.equal(readView(new URL("http://go/missing/path")).create,"missing");
 assert.equal(readView(new URL("http://go/?new=1")).create,"");
 assert.equal(readView(new URL("http://go/_/api/nope")).create,null);
 assert.equal(readView(new URL("http://go/_api")).create,"_api");
});
test("suggested names survive reload and malformed escapes do not crash the UI",()=>{
 assert.equal(readView(new URL("http://go/?new=1&short=missing")).create,"missing");
 assert.equal(readView(new URL("http://go/%FF")).create,"%FF");
});
test("failed requests retain server's useful error instead of reporting success",async()=>{
 const original=globalThis.fetch;
 globalThis.fetch=async()=>new Response(JSON.stringify({detail:"A link with that name already exists"}),{status:409,headers:{"Content-Type":"application/problem+json"}});
 try{await assert.rejects(()=>requestJSON("/_/api/links",{method:"POST",body:"{}"}),/already exists/);}finally{globalThis.fetch=original;}
});
test("mutation protection uses a header the browser permits",async()=>{
 const original=globalThis.fetch;
 globalThis.fetch=async(_input,init)=>{
  assert.equal(new Headers(init?.headers).get("X-Golink"),"1");
  return new Response(null,{status:204});
 };
 try{await requestJSON("/_/api/links/example",{method:"DELETE"});}finally{globalThis.fetch=original;}
});
