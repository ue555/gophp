<?php
$cache = new CachingIterator(new ArrayIterator(['a' => 10]), CachingIterator::FULL_CACHE);
$cache->rewind();
$copy = $cache->getCache();
$copy['a'] = 20;
$cache['b'] = 30;
echo $cache['a'], '|', $copy['a'], '|', isset($copy['b']) ? 'yes' : 'no', '|';
unset($cache['a']);
echo isset($cache['a']) ? 'yes' : 'no', '|', $copy['a'], '|', count($cache);
