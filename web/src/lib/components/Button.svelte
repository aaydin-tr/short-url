<script>
  import { Button } from "svelte-materialify";
  export let onClickHandler, text;

</script>

<Button {...$$restProps} on:click={onClickHandler}>{text}</Button>
