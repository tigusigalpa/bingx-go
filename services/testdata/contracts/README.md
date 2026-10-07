# Fixture provenance

`swap-market.documentation.json` records a review of pinned official references. It is documentation evidence, not a response from the exchange. The full docs-v3 bundle explicitly documents completed funding settlements and conflicts with the shorter AI reference's current/next description. Unknown trade pagination/quantity units and funding pagination completeness remain explicit. The funding response example contains only the first row from the official documentation example.

Files ending in `.synthetic.json` are invented loopback payloads for SDK boundary tests. Their contents do not establish exchange payload shape, chronology, settlement, quantity units, or candle closure. Tests elsewhere that generate responses in memory are synthetic as well.

No live captures are included. There are no credentials or production market calls in these tests.
