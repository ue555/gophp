<?php
$inner = new ArrayIterator(['a', 'b']);
$cache = new CachingIterator($inner, CachingIterator::FULL_CACHE);
$cache->rewind();
echo $cache->current(), '|', $inner->current(), '|', $cache->hasNext() ? 'more' : 'end', '|', count($cache);
$cache->next();
echo '|', $cache->current(), '|', $cache->hasNext() ? 'more' : 'end', '|', count($cache);
$cache->rewind();
echo '|', count($cache);
