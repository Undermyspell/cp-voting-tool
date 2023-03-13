<script lang="ts">
    import { LightSwitch, AppBar } from "@skeletonlabs/skeleton";
    import { userOnline } from "../lib/session";
</script>

<AppBar>
    <svelte:fragment slot="lead">
        <h3>Pollerr</h3>
    </svelte:fragment>
    <svelte:fragment slot="trail"
        ><div class="flex gap-4">
            <LightSwitch />
            <div>online: {$userOnline}</div>
        </div></svelte:fragment
    >
</AppBar>
