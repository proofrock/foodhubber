<script lang="javascript">
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
  // TODO i18n

  import "@materializecss/materialize/dist/css/materialize.min.css";

  import Header from "./components/app/Header.svelte";
  import { onDestroy, onMount } from "svelte";
  import { CALL, PAGES, ALERT_ERROR } from "./components/utils/Utils.svelte";
  import {
    SET_POLLING_CYCLE,
    START_GEN,
    STOP_GEN,
    SUB_WEEK_CHANGES,
  } from "./components/app/HubChecker.svelte";
  import LoginPage from "./components/LoginPage.svelte";
  import OrderPage from "./components/OrderPage.svelte";
  import ConsolePage from "./components/ConsolePage.svelte";
  import OrdersListPage from "./components/OrdersListPage.svelte";
  import StatsPage from "./components/StatsPage.svelte";
  import StockPage from "./components/StockPage.svelte";

  $: initData = null;
  $: checkout = null;
  $: week = 0;
  $: beneficiary = null;

  let page;

  $: {
    if (page !== undefined && page !== null)
      window.sessionStorage.setItem("page", page);
  }

  $: {
    if (!!initData) SET_POLLING_CYCLE(initData.polling_cycle);
  }

  onMount(async () => {
    const _page = window.sessionStorage.getItem("page");
    if (_page) page = parseInt(_page);
    else page = PAGES.ORDER;

    const _initData = window.sessionStorage.getItem("init_data");
    if (!window.sessionStorage.getItem("init_data")) {
      while (true) {
        let __initData = await CALL("getInitData");
        if (__initData.isErr) {
          await ALERT_ERROR(
            `<p>Init fallito.</p><p>${__initData.message}.</p>`,
          );
          continue;
        }
        __initData = __initData.payload;
        __initData.item_types = [];
        let curType = { color: "" };
        for (let i = 0; i < __initData.items.length; i++) {
          const item = __initData.items[i];
          if (item.color !== curType.color) {
            curType = {
              id: __initData.item_types.length + 1,
              color: item.color,
            };
            __initData.item_types.push(curType);
          }
          item.type = curType.id;
        }

        window.sessionStorage.setItem("init_data", JSON.stringify(__initData));
        initData = __initData;
        break;
      }
    } else initData = JSON.parse(_initData);

    initData.itemsMap = new Map();
    initData.items.forEach((it) => {
      initData.itemsMap.set(it.id, it);
    });

    initData.checkoutsMap = new Map();
    initData.checkouts.forEach((it) => {
      initData.checkoutsMap.set(it.id, it);
    });

    loadCheckout();

    SUB_WEEK_CHANGES(async (_week) => {
      week = _week;
    });
    START_GEN();
  });

  onDestroy(async () => {
    STOP_GEN(); // Really not needed
  });

  function doLogin(event) {
    checkout = event.detail;
    window.sessionStorage.setItem("checkout", JSON.stringify(checkout));
    loadCheckout();
  }

  function loadCheckout() {
    let _scheckout = window.sessionStorage.getItem("checkout");
    if (_scheckout !== null) {
      let _checkout = JSON.parse(_scheckout);
      checkout = initData.checkoutsMap.get(_checkout.id);
      checkout.cashier = _checkout.cashier;
    }
  }
</script>

<Header
  {initData}
  bind:beneficiary
  bind:checkout
  bind:week
  on:ch_page={(e) => {
    page = e.detail;
  }}
/>
<main>
  {#if initData !== null}
    {#if checkout === null}
      <LoginPage {initData} on:login={doLogin} />
    {:else if page === PAGES.ORDER}
      <OrderPage {initData} bind:beneficiary {checkout} />
    {:else if page === PAGES.ORDERS_LIST}
      <OrdersListPage {initData} />
    {:else if page === PAGES.STATS}
      <StatsPage {initData} />
    {:else if page === PAGES.STOCK}
      <StockPage {initData} />
    {:else if page === PAGES.CONSOLE}
      <ConsolePage {initData} />
    {/if}
  {/if}
</main>
