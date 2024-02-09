<script>
    /*
     * Copyright (C) 2024- Germano Rizzo
     *
     * This file is part of Sagramat.
     *
     * Sagramat is free software: you can redistribute it and/or modify
     * it under the terms of the GNU General Public License as published by
     * the Free Software Foundation, either version 3 of the License, or
     * (at your option) any later version.
     *
     * Sagramat is distributed in the hope that it will be useful,
     * but WITHOUT ANY WARRANTY; without even the implied warranty of
     * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     * GNU General Public License for more details.
     *
     * You should have received a copy of the GNU General Public License
     * along with Sagramat.  If not, see <http://www.gnu.org/licenses/>.
     */

    import { createEventDispatcher, onMount } from "svelte";
    import {
        ALERT_ERROR,
        MAT_ENABLE_SELECT,
    } from "../components/utils/Utils.svelte";
    import jsSHA from "jssha";

    const dispatch = createEventDispatcher();

    export let initData;

    let checkoutWithPassword = [];

    $: cashier = "";
    $: checkouts = [];
    $: selCheckout = "--";
    $: password = "";

    $: {
        const _checkouts = [];
        for (let i = 0; i < initData.checkouts.length; i++) {
            _checkouts.push(initData.checkouts[i].id);
            if (!!initData.checkouts[i].password)
                checkoutWithPassword.push(initData.checkouts[i].id);
        }
        checkouts = _checkouts;
        MAT_ENABLE_SELECT();
    }

    onMount(() => {
        MAT_ENABLE_SELECT();
    });

    async function log_in() {
        cashier = cashier.trim();
        password = password.trim();

        if (cashier === "") {
            await ALERT_ERROR("Devi indicare un operatore.");
            return;
        }

        if (selCheckout === "--") {
            await ALERT_ERROR("Devi indicare una postazione.");
            return;
        }

        const _checkout = initData.checkouts.find(
            (it) => it.id === selCheckout,
        );
        if (!_checkout) {
            await ALERT_ERROR("Postazione non valida.");
            return;
        }

        if (checkoutWithPassword.includes(selCheckout)) {
            // TODO the login should be made on the server and properly (401, cookie, ecc.)
            if (password === "") {
                await ALERT_ERROR(
                    "Devi inserire la password per questa postazione.",
                );
                return;
            }

            const shaObj = new jsSHA("SHA-256", "TEXT", { encoding: "UTF8" });
            shaObj.update(password);
            const pwdHash = shaObj.getHash("HEX");

            if (pwdHash !== _checkout.password) {
                await ALERT_ERROR("La password è errata.");
                return;
            }
        }

        // XXX We don't check if a checkout is already connected.

        const checkout = {
            id: selCheckout,
            cashier: cashier,
        };

        dispatch("login", checkout);
    }
</script>

<div class="row">&nbsp;</div>
<div class="row">
    <div class="col hide-on-small-and-down m3" />
    <div class="input-field col s12 m3 m-3">
        <i class="material-icons prefix">face</i>
        <input id="cashier" type="text" maxlength="16" bind:value={cashier} />
        <label for="cashier" class="active">Operatore</label>
    </div>
    <div class="input-field col s12 m3 m-3">
        <i class="material-icons prefix">toll</i>
        <select id="checkout" bind:value={selCheckout}>
            <option disabled selected>--</option>
            {#each checkouts as c}
                <option>{c}</option>
            {/each}
        </select>
        <label for="checkout">Postazione</label>
    </div>
    <div class="col hide-on-small-and-down m3" />
</div>
{#if checkoutWithPassword.includes(selCheckout)}
    <div class="row">
        <div class="col s1 m7" />
        <div class="input-field col s11 m2 m-3">
            <i class="material-icons prefix">lock</i>
            <input
                id="password"
                type="password"
                maxlength="32"
                bind:value={password}
            />
            <label for="password" class="active">Password</label>
        </div>
        <div class="col hide-on-small-and-down m3" />
    </div>
{/if}
<div class="center">
    <button class="waves-effect waves-light btn green" on:click={log_in}
        >Entra</button
    >
</div>
