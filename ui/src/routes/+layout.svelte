<script>
  import "../app.css";
  import { config } from "@fortawesome/fontawesome-svg-core";
  config.autoAddCss = false; // Tell Font Awesome to skip adding the CSS automatically since it's being imported
</script>

<slot />
