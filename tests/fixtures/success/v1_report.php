<?php
class Stats { public $accepted = 0; }
function summarize($rows, $stats) {
    $result = [];
    foreach ($rows as $key => $row) {
        if (!isset($row['name'], $row['amount']) || empty($row['name'])) continue;
        if (!is_numeric($row['amount'])) continue;
        $row['amount'] = (int)$row['amount'];
        unset($row['private']);
        $result[$key] = $row;
        ++$stats->accepted;
    }
    return $result;
}
$rows = ['a' => ['name' => 'tea', 'amount' => '12.9', 'private' => true], 'b' => ['name' => '0', 'amount' => 2], 'c' => ['name' => '水', 'amount' => 3], 'd' => ['name' => 'skip']];
$stats = new Stats; $result = summarize($rows, $stats);
foreach ($result as $key => $row) echo $key, ':', $row['name'], ':', $row['amount'], ':', strlen($row['name']), "\n";
echo count($result), ':', $stats->accepted, ':', isset($rows['a']['private']) + 0, ':', isset($result['a']['private']) + 0, ':', gettype($result['a']['amount']), "\n";
