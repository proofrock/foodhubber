<script>
    /*
     * Copyright (C) 2024- Germano Rizzo
     *
     * This file is part of FoodHubber.
     *
     * FoodHubber is free software: you can redistribute it and/or modify
     * it under the terms of the GNU General Public License as published by
     * the Free Software Foundation, either version 3 of the License, or
     * (at your option) any later version.
     *
     * FoodHubber is distributed in the hope that it will be useful,
     * but WITHOUT ANY WARRANTY; without even the implied warranty of
     * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     * GNU General Public License for more details.
     *
     * You should have received a copy of the GNU General Public License
     * along with FoodHubber.  If not, see <http://www.gnu.org/licenses/>.
     */

    import {
        SUB_ERROR_ON_SERVER,
        UNSUB_ERROR_ON_SERVER,
    } from "./HubChecker.svelte";
    import { onDestroy, onMount } from "svelte";

    let id = "" + Math.random();

    $: errorServer = false;

    onMount(async () => {
        SUB_ERROR_ON_SERVER(id, (bool) => {
            errorServer = bool;
        });
    });

    onDestroy(() => {
        UNSUB_ERROR_ON_SERVER();
    });
</script>

{#if errorServer}
    <span class="blink red-text">
        <i class="material-icons" title="Non connesso!">portable_wifi_off</i>
    </span>
{/if}

<style>
    .blink {
        -webkit-animation: blink 1s step-end infinite;
        -moz-animation: blink 1s step-end infinite;
        -o-animation: blink 1s step-end infinite;
        animation: blink 1s step-end infinite;
    }

    @-webkit-keyframes blink {
        67% {
            opacity: 0;
        }
    }

    @-moz-keyframes blink {
        67% {
            opacity: 0;
        }
    }

    @-o-keyframes blink {
        67% {
            opacity: 0;
        }
    }

    @keyframes blink {
        67% {
            opacity: 0;
        }
    }
</style>
