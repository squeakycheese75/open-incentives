# Open Incentives

Open Incentives is an open-source promotion engine for modern commerce.

Build discounts, coupons, loyalty rewards, referral campaigns, and other
incentives through a lightweight, API-first platform. If you're
evaluating platforms like Talon.One or Voucherify, Open Incentives
offers a self-hosted, open-source alternative.

## Features

-   API-first promotion engine
-   Self-hosted with Docker
-   Open source
-   Lightweight and developer-friendly
-   Rule-based campaigns
-   Discounts, coupons, and incentives
-   Built for modern commerce

## Quick Start

``` bash
git clone https://github.com/squeakycheese75/open-incentives.git
cd open-incentives
docker compose up 
```

or with demo store 
``` bash
git clone https://github.com/squeakycheese75/open-incentives.git
cd open-incentives
docker compose up --profile demo-store
```


Once running:

-   API: http://localhost:8080
-   Admin Portal: http://localhost:3001
-   Demo Store: http://localhost:3000 

Log in to the Admin Portal with the default bootstrapped credentials:

-   Email: `admin@example.com`
-   Password: `change-me`

(Set `BOOTSTRAP_ADMIN_EMAIL` / `BOOTSTRAP_ADMIN_PASSWORD` before first run to change these.)

See the `docs/Quickstart.md` guide for a complete walkthrough.

## Why?

Almost every commerce business eventually needs promotions:

-   Coupon codes
-   Discounts
-   Store credit
-   Loyalty rewards
-   Referral campaigns
-   Win-back offers

Today, the options are often:

1.  Build and maintain a custom promotion engine.
2.  Buy an enterprise incentives platform like Talon.One or Voucherify.
3.  Delay experimentation because the investment is too high.

Open Incentives fills the gap between basic discount functionality and
enterprise promotion platforms.

Our goal is to make promotions:

-   Easy to build
-   Easy to integrate
-   Easy to operate
-   Easy to remove

Open Incentives is API-first, self-hosted, and open source---giving
teams full control over their promotion infrastructure without the
complexity or cost of enterprise platforms.

## Roadmap

-   [x] Rule evaluation engine
-   [x] Campaign management API
-   [x] API key authentication
-   [x] Demo store UI
-   [x] Admin UI
-   [ ] PostgreSQL storage
-   [ ] Webhooks
-   [ ] Loyalty points
-   [ ] A/B testing
