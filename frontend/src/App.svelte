<script lang="ts">
 import {onMount,tick} from "svelte";
 import {Button,SearchInput,SelectDropdown,ThemeToggle,EmptyState} from "@kenn-io/kit-ui";
 import LinkList from "$lib/components/LinkList.svelte";
 import LinkEditor from "$lib/components/LinkEditor.svelte";
 import TemplateHelp from "$lib/components/TemplateHelp.svelte";
 import LinkImport from "$lib/components/LinkImport.svelte";
 import {fetchLinks,fetchSettings,filterLinks,sortLinks,readView,type Link,type Settings,type View,type ImportResult} from "$lib/links.js";
 let links=$state<Link[]>([]);let settings=$state<Settings|null>(null);let loaded=$state(false);let loading=$state(false);let error=$state("");let status=$state("");
 const initialView=readView(new URL(window.location.href));
 let view=$state(initialView);let query=$state(initialView.query);let sort=$state<"name"|"clicks">("clicks");let request:AbortController|null=null;

 let importButton:HTMLButtonElement|null=$state(null);let importing=$state(false);let importStatus=$state("");
 const visible=$derived(sortLinks(filterLinks(links,query),sort));
 const normalize=(value:string)=>value.toLowerCase().replaceAll("-","");
 const editing=$derived(view.edit===null?null:links.find(link=>normalize(link.short)===normalize(view.edit!))??null);
 const editorOpen=$derived(settings!==null&&!settings.readOnly&&(view.create!==null||editing!==null));

 function setView(patch:Partial<View>,replace=false){
  if(patch.edit!==undefined||patch.create!==undefined)importing=false;
  view={...view,...patch};query=view.query;const url=new URL("/",window.location.origin);
  if(view.query)url.searchParams.set("q",view.query);if(view.edit!==null)url.searchParams.set("edit",view.edit);
  if(view.create!==null){url.searchParams.set("new","1");if(view.create)url.searchParams.set("short",view.create);}
  if(view.help)url.searchParams.set("help","1");
  window.history[replace?"replaceState":"pushState"]({},"",url);
 }
 async function load(){if(loading)return;loading=true;error="";const controller=new AbortController();request=controller;
  try{const [catalog,config]=await Promise.all([fetchLinks(controller.signal),fetchSettings(controller.signal)]);links=catalog.links;settings=config;loaded=true;}
  catch(e){if(!controller.signal.aborted)error=loaded?"Couldn't refresh links. Showing the last loaded list.":e instanceof Error?e.message:"Couldn't load links. Try again.";}
  finally{if(!controller.signal.aborted)loading=false;}
 }
 async function closeEditor(focusShort?:string){setView({edit:null,create:null});await tick();const target=focusShort?document.getElementById(`edit-${focusShort}`):null;(target?.querySelector<HTMLButtonElement>("button")??document.querySelector<HTMLButtonElement>(".new-link"))?.focus();}
 function saved(link:Link){links=[...links.filter(item=>normalize(item.short)!==normalize(link.short)),link];status=`${link.short} saved.`;void closeEditor(link.short);}
 function deleted(short:string){links=links.filter(link=>link.short!==short);status=`${short} deleted.`;void closeEditor();}
 function openImport(){setView({edit:null,create:null});importStatus="";importing=true;}
 async function closeImport(reload=false){importing=false;if(reload)await load();await tick();importButton?.focus();}
 function imported(result:ImportResult){importStatus=`Imported ${result.imported} ${result.imported===1?"link":"links"}. Skipped ${result.skipped} existing ${result.skipped===1?"link":"links"}.`;if(result.skippedStats)importStatus+=` Skipped ${result.skippedStats} history rows for links outside the backup.`;void closeImport(true);}
 onMount(()=>{void load();const onPop=()=>{importing=false;view=readView(new URL(window.location.href));query=view.query;};window.addEventListener("popstate",onPop);return()=>{request?.abort();window.removeEventListener("popstate",onPop);};});
</script>

<main class="app-shell">
 <header class="app-header"><a class="app-title" href="/" onclick={(event)=>{event.preventDefault();setView({query:"",edit:null,create:null,help:false});}}><svg aria-hidden="true" viewBox="0 0 32 32" fill="none" stroke="currentColor" stroke-width="2"><path d="m13 19 6-6M11.5 21.5l-2 2a5.66 5.66 0 0 1-8-8l6-6a5.66 5.66 0 0 1 8 0m5 11a5.66 5.66 0 0 0 8 0l6-6a5.66 5.66 0 0 0-8-8l-2 2" transform="translate(2 1) scale(.82)" stroke-linecap="round" stroke-linejoin="round"/></svg><h1>Golinks</h1></a>
  <div class="theme-control"><ThemeToggle/></div>
 </header>
 {#if error}<div class="notice" role="alert"><p>{error}</p><Button class="control" onclick={load} disabled={loading}>Try again</Button></div>{/if}
 <p class="sr-only" role="status" aria-live="polite">{status}</p>
 {#if importStatus}<p class="import-status" role="status">{importStatus}</p>{/if}
 {#if importing&&settings&&!settings.readOnly}<LinkImport onImport={imported} onCancel={()=>void closeImport()}/>{/if}
 {#if editorOpen}{#key editing?.short??`new-${view.create}`}<LinkEditor link={editing} suggested={view.create??""} onSave={saved} onDelete={deleted} onCancel={()=>void closeEditor()}/>{/key}
 {:else if loaded&&view.edit!==null&&!editing}<p class="notice" role="alert">That link no longer exists. <button class="text-link" onclick={()=>void closeEditor()}>Back to links</button></p>{/if}
 <section class="directory" aria-label="Links" aria-busy={loading}>
  <div class="directory-toolbar"><div class="search-field"><SearchInput id="search" block ariaLabel="Search links" value={query} oninput={(value)=>setView({query:value},true)} placeholder="Search links"/></div>
   {#if !settings?.readOnly}<Button class="control new-link" onclick={()=>{status="";setView({edit:null,create:""});}} disabled={!settings||editorOpen||importing}><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M12 5v14M5 12h14" stroke-linecap="round"/></svg>New link</Button>{/if}
  {#if settings?.readOnly}<span>Read-only</span>{/if}<div class="sort-label"><span>Sort</span><SelectDropdown title="Sort links" value={sort} onchange={(value)=>sort=value==="clicks"?"clicks":"name"} options={[{value:"name",label:"Name"},{value:"clicks",label:"Most used"}]} align="end"/></div></div>
  {#if !loaded&&loading}<div class="empty-state" role="status"><p>Loading links…</p></div>
  {:else if loaded&&links.length===0}<div class="empty-state"><EmptyState title="No links yet" description={settings?.readOnly?"No links have been added to this directory.":"Create a short name for a URL you use often."}/></div>
  {:else if loaded&&visible.length===0}<div class="empty-state"><EmptyState title="No matching links" description="Try another name or destination."><Button onclick={()=>setView({query:""},true)}>Clear search</Button></EmptyState></div>
  {:else if loaded&&settings}<LinkList links={visible} publicUrl={settings.publicUrl} readOnly={settings.readOnly} disabled={importing} onEdit={(link)=>{status="";setView({edit:link.short,create:null});}}/>{/if}
 </section>
 <footer class="app-footer"><button class="footer-action" onclick={()=>setView({help:!view.help})} aria-expanded={view.help}>Template help</button><a href="/_/export">Export links</a>{#if settings&&!settings.readOnly}<button class="footer-action" bind:this={importButton} onclick={openImport} aria-expanded={importing} disabled={importing||editorOpen||loading}>Import links</button>{/if}<a href="/_/api/docs">API docs</a></footer>
 {#if view.help}<TemplateHelp/>{/if}
</main>
