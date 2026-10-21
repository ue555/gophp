<?php
$cache = new CachingIterator(new ArrayIterator([1]), 65536);
echo $cache->getFlags(), ':', $cache->valid() ? 'valid' : 'end', '|';
$cache->setFlags(65536);
echo $cache->getFlags(), '|';
$cache->rewind();
echo $cache->getFlags(), ':', $cache->valid() ? 'valid' : 'end', '|';
$cache->setFlags(0);
echo $cache->getFlags(), '|';
$cache->setFlags(65536 | 512);
echo $cache->getFlags(), '|';
$other = new CachingIterator(new ArrayIterator([]), 512);
echo $other->getFlags();
