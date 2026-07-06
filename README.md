# Open Incentives

> Open-source promotions infrastructure for teams that want to experiment with incentives without making a long-term ecosystem commitment.

## Why?

Many companies eventually need promotions:

- Coupon codes
- Discounts
- Credits
- Basic loyalty
- Referral rewards
- Win-back campaigns

Today the choices are often:

1. Build it yourself.
2. Buy an enterprise incentives platform.
3. Delay experimentation.

Open Incentives aims to fill the gap between simple discount functionality and full enterprise incentives suites.

The goal is to make promotions:

- Easy to add
- Easy to integrate
- Easy to remove

---

# Vision

Become the default open-source promotions engine that teams reach for when they need incentives quickly and don't want to commit to an entire ecosystem.

The project is designed for:

- Custom ecommerce applications
- SaaS products
- Marketplaces
- Platforms embedding promotions into their products

---

# Philosophy

## Start simple

A team should be able to:

1. Start the service
2. Create a campaign
3. Integrate with their application
4. See a promotion working

...in less than an afternoon.

## Avoid lock-in

Promotions should not require a strategic vendor decision.

The system should be:

- Self-hostable
- Portable
- Easy to uninstall
- Easy to migrate away from

## Build from the core

The heart of the project is a standalone rules engine.

Facts + Rules → Promotion Engine → Actions + Decisions

Everything else is built around this core.

---

# Initial Use Case

A company running on:

Frontend → Backend → Stripe → No promotions infrastructure

wants to launch:

> 10% off orders over €50

without building a promotions system from scratch or adopting an enterprise platform.

---

# Project Goals

## V1

- Rules engine
- Promotion evaluation API
- Campaign management
- Coupons
- Event tracking
- Basic admin UI
- Demo application

## Future

- Loyalty
- Referrals
- Segmentation
- Experimentation
- Analytics
- Multi-tenancy

---

# Non-Goals

Open Incentives is not trying to be:

- A CRM
- A CDP
- Marketing automation software
- An email platform
- A customer journey tool
- An enterprise incentives suite (yet)

---

# Architecture

Admin UI
↑
API + Persistence
↑
Promotion Engine

The engine is intentionally designed to be:

- Deterministic
- Importable
- Portable
- Versioned
- Usable as a standalone dependency

---

# Guiding Principles

## Easy to Add

docker compose up

## Easy to Integrate

result, err := client.Evaluate(...)

## Easy to Remove

Delete a few API calls and shut down the service.

No lock-in.

---

# Status

🚧 Early development.

The project is currently focused on validating the core promotion engine and the first end-to-end use cases.
