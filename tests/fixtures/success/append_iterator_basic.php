<?php
$combined = new AppendIterator;
$combined->append(new ArrayIterator(['a' => 1]));
$combined->append(new ArrayIterator([]));
$combined->append(new ArrayIterator(['a' => 2, 3]));
foreach ($combined as $key => $value) {
    echo $combined->getIteratorIndex(), ':', $key, '=', $value, ';';
}
echo '|', count($combined->getArrayIterator());
