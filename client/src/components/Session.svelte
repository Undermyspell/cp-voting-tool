<script lang="ts">
    import { ArrowUpCircle, Play, Stop } from "@steeze-ui/remix-icons";
    import { Icon } from "@steeze-ui/svelte-icon";
    import { startSession, stopSession } from "../lib/session";
    import { ButtonGroup, Button } from "flowbite-svelte";
</script>

<ButtonGroup>
    <Button type="button" on:click={startSession}>
        <Icon src={Play} size="20" />
        <span class="ml-2">Fragerunde starten</span>
    </Button>
    <Button type="button" on:click={stopSession}>
        <Icon src={Stop} size="20" />
        <span class="ml-2"> Fragerunde beenden </span>
    </Button>
</ButtonGroup>
