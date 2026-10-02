-- stats: raw aggregates; R/win-rate/expectancy/drawdown math lives in pkg/calc (unit tested)

-- name: StatsClosedTrades :many
select id, pnl, r_multiple, followed_rules, opened_at, closed_at, pair, setup_id, account_id
from trades
where user_id = $1 and status = 'closed'
  and (sqlc.narg(account_id)::uuid is null or account_id = sqlc.narg(account_id)::uuid)
order by closed_at asc nulls last;

-- name: StatsEquityCurve :many
select closed_at, pnl from trades
where user_id = $1 and status = 'closed'
  and (sqlc.narg(account_id)::uuid is null or account_id = sqlc.narg(account_id)::uuid)
order by closed_at asc nulls last;

-- name: StatsBySetup :many
select setup_id, count(*) as trades, sum(pnl) as total_pnl, avg(r_multiple) as avg_r
from trades where user_id = $1 and status = 'closed'
  and (sqlc.narg(account_id)::uuid is null or account_id = sqlc.narg(account_id)::uuid)
group by setup_id;

-- name: StatsByPair :many
select pair, count(*) as trades, sum(pnl) as total_pnl, avg(r_multiple) as avg_r
from trades where user_id = $1 and status = 'closed'
  and (sqlc.narg(account_id)::uuid is null or account_id = sqlc.narg(account_id)::uuid)
group by pair;
