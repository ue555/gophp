<?php
$cache = new CachingIterator(new ArrayIterator(['a' => 'x', 'b' => 'y', 'c' => 'z']));
foreach ($cache as $key => $value) {
    echo $key, '=', $value, ':', $cache->hasNext() ? 'more' : 'end', ';';
}
