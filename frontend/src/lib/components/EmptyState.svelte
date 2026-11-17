<script lang="ts">
  import type { Snippet } from "svelte";

  let { title, description, children }: { title: string; description: string; children?: Snippet } =
    $props();
</script>

<div class="empty-state">
  <h2>{title}</h2>
  <p>{description}</p>
  {#if children}<div class="empty-state-action">{@render children()}</div>{/if}
</div>
