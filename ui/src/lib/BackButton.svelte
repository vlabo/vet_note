<script>
  import { faAngleLeft } from "@fortawesome/free-solid-svg-icons";
  import Icon from "$lib/icon.svelte";

  function goBack() {
    history.back();
  }
</script>

<button
  class="text-black hover:bg-blue-100 hover:rounded-lg focus:outline-none py-2 px-4"
  on:click={goBack}
>
  <Icon icon={faAngleLeft} />
</button>
