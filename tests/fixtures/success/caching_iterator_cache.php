<?php
$cache = new CachingIterator(new ArrayIterator(['x' => 1, 'y' => 2]), CachingIterator::FULL_CACHE);
foreach ($cache as $value) {
    echo count($cache), ':';
}
echo '|', count($cache->getCache()), '|', $cache['x'];
