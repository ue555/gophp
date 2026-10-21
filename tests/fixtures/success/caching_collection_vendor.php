<?php
require 'research/laravel-api/vendor/autoload.php';
$collection = new \Illuminate\Support\Collection(['x', 'y', 'z']);
$plain = $collection->getCachingIterator();
foreach ($plain as $value) {
    echo $value;
}
echo '|';
$cache = $collection->getCachingIterator(CachingIterator::FULL_CACHE);
foreach ($cache as $value) {
    echo $value, $cache->hasNext() ? ',' : '.';
}
echo '|', count($cache->getCache());
