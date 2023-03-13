<script lang="ts">
    import { LightSwitch, AppBar } from "@skeletonlabs/skeleton";
</script>

<AppBar>
    <svelte:fragment slot="lead">
        <h3>Pollerr</h3>
    </svelte:fragment>
    <svelte:fragment slot="trail"><LightSwitch /></svelte:fragment>
</AppBar>
