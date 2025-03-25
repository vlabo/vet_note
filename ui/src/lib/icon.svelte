<script lang="ts">
  import { FontAwesomeIcon } from "@fortawesome/svelte-fontawesome";

  import "@fortawesome/fontawesome-svg-core/styles.css";
  export let icon: any;
  let className = "";
  export { className as class };
</script>

<FontAwesomeIcon {icon} class={className} />
