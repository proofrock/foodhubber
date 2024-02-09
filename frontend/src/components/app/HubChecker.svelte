<script context="module">
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

    import { LOGOUT } from "../utils/Utils.svelte";
    import { CALL } from "../utils/Utils.svelte";

    let CYCLE_PERIOD = 2000;

    let lambdaForOrders = null;
    let lambdaForStock = null;
    let lambdaForWeek = null;

    let cycleHandler = null;

    let lastRunId = null;
    let lastOrdersGen = null;
    let lastStockGen = null;
    let lastWeek = 0;

    let lambdaForErrorServer = {};
    let errorServer = false;

    // Returns if a logout was performed
    function ckRunId(gen) {
        if (lastRunId === null) lastRunId = gen.run_id;
        if (lastRunId !== gen.run_id) {
            LOGOUT();
            return true;
        }
        return false;
    }

    async function ckOrdersGen(gen) {
        if (lastOrdersGen !== gen) {
            if (lastOrdersGen !== null && lambdaForOrders !== null)
                await lambdaForOrders();
            lastOrdersGen = gen;
        }
    }

    async function ckStockGen(gen) {
        if (lastStockGen !== gen) {
            if (lastStockGen !== null && lambdaForStock !== null)
                await lambdaForStock();
            lastStockGen = gen;
        }
    }

    async function ckWeek(week) {
        if (lastWeek !== week) {
            if (lambdaForWeek !== null && lambdaForWeek !== null)
                await lambdaForWeek(week);
            lastWeek = week;
        }
    }

    async function doCycle() {
        const checkout = window.sessionStorage.getItem("checkout");

        const map = {};
        if (!!checkout) {
            const _checkout = JSON.parse(checkout);
            map["pos"] = _checkout.id;
            map["op"] = _checkout.cashier;
        }
        const res = await CALL("doCycle", "POST", null, map);
        signalErrorOnServer(res.isErr);
        if (!res.isErr) {
            const ret = res.payload;

            if (ckRunId(ret)) return;

            await ckWeek(ret.week);
            await ckOrdersGen(ret.gen_orders);
            await ckStockGen(ret.gen_stock);
        }
    }

    export const SET_POLLING_CYCLE = function (pollingCycle) {
        CYCLE_PERIOD = pollingCycle;
    };

    export const SUB_ORDERS_CHANGES = function (_lambda) {
        lambdaForOrders = _lambda;
    };

    export const SUB_STOCK_CHANGES = function (_lambda) {
        lambdaForStock = _lambda;
    };

    export const SUB_WEEK_CHANGES = function (_lambda) {
        lambdaForWeek = _lambda;
    };

    // XXX Eventually: verify if this happen strictly before a new subscription (onDestroy of old page happens before
    //      onMount of a new one). If not, Bad Things Will Happen.
    //      Preliminary tests show that it behaves correctly. I'll leave it here for future reference.
    export const UNSUB_ORDERS_CHANGES = function () {
        lambdaForOrders = null;
    };

    export const UNSUB_STOCK_CHANGES = function () {
        lambdaForStock = null;
    };

    export const SUB_ERROR_ON_SERVER = function (id, _lambda) {
        lambdaForErrorServer[id] = _lambda;
    };

    function signalErrorOnServer(bool) {
        if (bool !== errorServer) {
            errorServer = bool;
            Object.keys(lambdaForErrorServer).forEach((key) =>
                lambdaForErrorServer[key](errorServer),
            );
        }
    }

    export const UNSUB_ERROR_ON_SERVER = function (id) {
        delete lambdaForErrorServer[id];
    };

    async function cycle() {
        try {
            await doCycle();
        } catch (e) {
            console.error(e);
        }
        cycleHandler = setTimeout(cycle, CYCLE_PERIOD);
    }

    export const START_GEN = function () {
        cycleHandler = setTimeout(cycle, CYCLE_PERIOD);
    };

    export const STOP_GEN = function () {
        clearTimeout(cycleHandler);
    };
</script>
