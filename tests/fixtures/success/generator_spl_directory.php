<?php
$base = __DIR__ . '/../support/spl_walk';
$outer = new RecursiveDirectoryIterator($base, RecursiveDirectoryIterator::SKIP_DOTS);
$walk = new RecursiveIteratorIterator($outer, RecursiveIteratorIterator::LEAVES_ONLY);
$a = 0;
$b = 0;
$depth = 0;
foreach ($walk as $file) {
    if ($file->getFilename() === 'a.txt') $a++;
    if ($file->getFilename() === 'b.txt') { $b++; $depth = $walk->getDepth(); }
}
echo $a, '|', $b, '|', $depth;
