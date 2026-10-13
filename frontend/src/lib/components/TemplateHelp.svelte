<script lang="ts">
 const helpers=["PathEscape","QueryEscape","TrimPrefix","TrimSuffix","ToLower","ToUpper","Match"];
 const example='https://www.google.com/{{if .Path}}search?q={{QueryEscape .Path}}{{end}}';
</script>
<section class="template-help" aria-labelledby="help-title">
 <h2 id="help-title">Destination templates</h2>
 <p>A plain URL appends the extra path. Go templates let you put that path anywhere in the destination.</p>
 <pre><code>{example}</code></pre>
 <p>With this destination, <code>go/search/hello world</code> opens a Google search. <code>go/search</code> opens Google.</p>
 <dl><div><dt><code>.Path</code></dt><dd>The path after the short name, without its leading slash.</dd></div><div><dt><code>.Now</code></dt><dd>The current UTC time. Use <code>{'{{.Now.Format "2006-01-02"}}'}</code> for today's date.</dd></div></dl>
 <p class="helper-label">Available helpers</p><div class="helper-list">{#each helpers as helper}<code>{helper}</code>{/each}</div>
 <p>Query parameters from the incoming link are added to the destination. Relative URLs can point to another short link.</p>
</section>
