<?php
$cache = new CachingIterator(new ArrayIterator(['a' => 10]), CachingIterator::FULL_CACHE);
$cache->rewind();
echo $cache->offsetExists('a') ? 'yes|' : 'no|';
echo $cache->offsetGet('a'), '|';
$cache->offsetSet('b', 20);
echo count($cache), ':', $cache->offsetGet('b'), '|';
$cache->offsetUnset('a');
echo $cache->offsetExists('a') ? 'yes|' : 'no|';
set_error_handler(function ($level, $message) {
    echo $level, ':', $message, '|';
    return true;
});
echo $cache->offsetGet('missing') === null ? 'null' : 'wrong';
restore_error_handler();
