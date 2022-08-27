<script>
  import { CardText } from "svelte-materialify";
  export let message, short_url;
</script>

<CardText>
  {message}
  <br />
  <a href={short_url} target="_blank">{short_url}</a>
</CardText>
