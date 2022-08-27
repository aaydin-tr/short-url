<script>
  import { TextField } from "svelte-materialify";
  export let error, messages, url, onKeyPress;
</script>

<TextField
  {...$$restProps}
  messages={(error && messages) || []}
  success={!error}
  {error}
  bind:value={url}
  on:keypress={onKeyPress}
/>
