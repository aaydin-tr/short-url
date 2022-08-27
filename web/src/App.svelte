<script>
  import { Card, CardText, CardActions, MaterialApp } from "svelte-materialify";
  import Link from "./lib/components/Link.svelte";
  import Button from "./lib/components/Button.svelte";
  import Input from "./lib/components/Input.svelte";
  import { handleNewShortURL } from "./lib/util/util";
  let url, error, messages, response;


  const onKeyPress = (e) => {
    if (e.charCode === 13) handleNewShortURL();
  };

  const onClickHandler = async (url) => {
    [response, error, messages] = await handleNewShortURL(url);
  };
</script>

<div class="main">
  <div class="app">
    <MaterialApp theme="dark">
      <div class="d-flex justify-center mt-4 mb-4">
        <Card rounded outlined style="width:500px;">
          <div class="pl-4 pr-4 pt-3">
            <span class="text-h5 mb-2">Short URL</span>
            <br />
          </div>
          <CardText>Create New Short Url</CardText>
          <CardText>
            <Input
              {error}
              {messages}
              {onKeyPress}
              bind:url
              type="url"
              name="url"
              id="url"
              placeholder="https://example.com"
              pattern="https://.*"
              autocomplete="off"
              clearable={true}
            /></CardText
          >
          {#if !error && response?.data}
            <Link message={messages} short_url={response.data?.short_url} />
          {/if}
          <CardActions>
            <Button
              text={"Create"}
              onClickHandler={() => onClickHandler(url)}
              outlined
              rounded
              size="small"
              class="blue blue-text"
              style="margin: 10px auto 0px;"
            />
          </CardActions>
        </Card>
      </div>
    </MaterialApp>
  </div>
</div>

<style>
  .app {
    margin: auto;
    width: auto;
  }
</style>
