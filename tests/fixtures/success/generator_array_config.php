<?php
$data = ['b' => 2, 'a' => 1];
ksort($data, SORT_NATURAL);
$keys = array_keys($data);
$mapped = array_map(fn ($v, $k) => $v . $k, $data, $keys);
$result = array_combine($keys, $mapped);
echo implode(',', array_keys($result)), '|', $result['a'], '|', $result['b'], '|';
$removed = array_splice($keys, 0, 1);
echo $removed[0], '|', $keys[0];
