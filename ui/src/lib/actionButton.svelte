<script lang="ts">
  import Icon from "$lib/icon.svelte";
  import { faPlus } from "@fortawesome/free-solid-svg-icons";

  var hrefProp = "";
  export { hrefProp as href };
</script>

<a
  class="fixed bottom-6 right-6 bg-blue-600 hover:bg-blue-700 transition text-white w-14 h-14 rounded-full flex items-center justify-center shadow-lg focus:outline-none"
  href={hrefProp}
>
  <Icon icon={faPlus} />
</a>
