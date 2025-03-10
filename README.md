# citylyf

![Screenshot](/screenshot.png)

`citylyf` is an attempt at a very simple economic/governance sim written in Go where the player runs a small city state.

Currently, the player can create houses where people can move in. There are companies at which people can
get jobs. They pay taxes and rent. The interest rate is set by the Central Bank to counter inflation.

## Planned Todos

- [x] Turn people, households and companies into a map
- [ ] Household Budgeting - think about childcare expenses, groceries, shopping, vacation, utilities etc
- [x] Housing market - rent, no. of bedrooms etc., grow rent yearly by inflation rate
- [ ] People should marry, have babies, get promoted, move out out the house, die etc.
- [ ] Yearly budget - once a year, we show users government income vs expenditure and store these values for recall
- [ ] Calculate realistic government expenses - e.g. laying down roads and building houses should cost the govt money
- [ ] Pension fund with employee + employer + government contributions
- [ ] Companies should be tied to office space/industrial space availability
- [x] Retail companies + shops
- [ ] Shops/offices should be like houses, built and kept unoccupied until a company moves in
- [ ] Companies with no employees for a year should shut down (tie productivity to employee count?)
- [ ] Land use type to track tile land use instead of the current booleans
- [ ] Regions to track population and simulate traffic between them
- [ ] Forests and farmland
- [x] Housing estates instead of laying down individual houses?
- [x] Better UI for road laying
- [x] To enable both of the above, add a way to create rectangular plots using the mouse
- [ ] "Buildable" bool on tiles
- [ ] Build animation
- [x] Break sprites into multiple files
- [x] Save and load games