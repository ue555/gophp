<?php
class WalkStore {
    public $rows = [['a' => 1, 'b' => 2]];
    public $tree = ['child' => ['x' => 3, 'y' => 4]];
}
$store = new WalkStore;
$receiverCalls = 0;
$keyCalls = 0;
function walkStore() { global $store, $receiverCalls; $receiverCalls++; return $store; }
function walkIndex() { global $keyCalls; $keyCalls++; return 0; }
echo array_walk(callback: function (&$value) { $value += 5; }, array: walkStore()->rows[walkIndex()]) ? 'ok|' : 'bad|';
echo $receiverCalls, '|', $keyCalls, '|', $store->rows[0]['a'], '|', $store->rows[0]['b'], '|';
echo array_walk_recursive(callback: function (&$value) { $value *= 2; }, array: walkStore()->tree) ? 'ok|' : 'bad|';
echo $receiverCalls, '|', $store->tree['child']['x'], '|', $store->tree['child']['y'];
