export interface Link {short:string;long:string;created:string;lastEdit:string;clicks:number}
export interface Settings {publicUrl:string;readOnly:boolean}
export interface ImportResult {imported:number;skipped:number;skippedStats:number}
export interface View {query:string;edit:string|null;create:string|null;help:boolean}

export function readView(url:URL):View {
 const params=url.searchParams;
 let unknown:string|null=null;
 if(url.pathname!=="/"&&!url.pathname.startsWith("/.")&&url.pathname!=="/_"&&!url.pathname.startsWith("/_/")){
  const segment=url.pathname.split("/")[1]??"";
  try{unknown=decodeURIComponent(segment);}catch{unknown=segment;}
 }
 return {query:params.get("q")??"",edit:params.get("edit"),create:params.has("new")?params.get("short")??"":unknown,help:params.has("help")};
}
export function filterLinks(links:Link[],query:string):Link[]{
 const needle=query.trim().toLowerCase();return links.filter(link=>`${link.short} ${link.long}`.toLowerCase().includes(needle));
}
export function sortLinks(links:Link[],sort:"name"|"clicks"):Link[]{
 return [...links].sort((a,b)=>(sort==="clicks"?b.clicks-a.clicks:0)||a.short.localeCompare(b.short));
}
export async function requestJSON<T>(path:string,init:RequestInit={}):Promise<T>{
 const headers=new Headers(init.headers);
 if(init.method && init.method!=="GET"){headers.set("Content-Type","application/json");headers.set("X-Golink","1");}
 const response=await fetch(path,{...init,headers});
 if(!response.ok){let message=`Request failed (${response.status}). Try again.`;
  const text=await response.text();try{const error=JSON.parse(text) as {detail?:string;errors?:{message?:string}[]};message=error.detail??error.errors?.[0]?.message??message;}catch{if(text.trim())message=text.trim();}
  throw new Error(message);
 }
 return response.status===204?undefined as T:await response.json() as T;
}
export const fetchLinks=(signal?:AbortSignal)=>requestJSON<{links:Link[]}>("/_/api/links",{signal});
export const fetchSettings=(signal?:AbortSignal)=>requestJSON<Settings>("/_/api/settings",{signal});
export const fetchLink=(short:string,signal?:AbortSignal)=>requestJSON<Link>(`/_/api/links/${encodeURIComponent(short)}`,{signal});
export const createLink=(short:string,long:string,signal?:AbortSignal)=>requestJSON<Link>("/_/api/links",{method:"POST",body:JSON.stringify({short,long}),signal});
export const updateLink=(short:string,long:string,signal?:AbortSignal)=>requestJSON<Link>(`/_/api/links/${encodeURIComponent(short)}`,{method:"PUT",body:JSON.stringify({long}),signal});
export const deleteLink=(short:string,signal?:AbortSignal)=>requestJSON<void>(`/_/api/links/${encodeURIComponent(short)}`,{method:"DELETE",signal});
export const previewLink=(long:string,path:string,query:string,signal?:AbortSignal)=>requestJSON<{url:string}>("/_/api/resolve",{method:"POST",body:JSON.stringify({long,path,query}),signal});
export const importLinks=(links:string,stats?:string)=>requestJSON<ImportResult>("/_/api/import",{method:"POST",body:JSON.stringify({links,stats})});
