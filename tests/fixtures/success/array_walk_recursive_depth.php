<?php
$tree = 7;
for ($i = 0; $i < 70; $i++) {
    $tree = [$tree];
}

$calls = 0;
array_walk_recursive($tree, function (&$value) use (&$calls) {
    $calls++;
    $value++;
});

for ($i = 0; $i < 70; $i++) {
    $tree = $tree[0];
}
echo $tree, '|', $calls;
