<?php
$cache = new CachingIterator(new ArrayIterator(['a', 'b']), CachingIterator::FULL_CACHE);
$cache->rewind();
echo $cache->getFlags(), ':', count($cache), '|';
$cache->setFlags(0);
echo $cache->getFlags(), '|';
$cache->setFlags(CachingIterator::FULL_CACHE);
echo count($cache), '|';
$cache->next();
echo $cache->current(), ':', count($cache), '|';
$cache->next();
echo $cache->valid() ? 'valid' : 'end', ':', $cache->getFlags();
