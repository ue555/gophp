<?php
trait ArtisanProbeTrait {}
class ArtisanProbeParent { use ArtisanProbeTrait; }
class ArtisanProbeChild extends ArtisanProbeParent {}
echo sprintf('%2$s|%1$s', 'a', 'b'), '|';
echo implode(',', array_intersect([1, 2, 3], [2, 3], [3])), '|';
echo implode(',', class_parents(ArtisanProbeChild::class)), '|';
echo implode(',', class_uses(ArtisanProbeParent::class)), '|';
$info = new SplFileInfo(__DIR__ . '/../support/spl_walk/a.txt');
echo $info->getBasename('.txt'), ':', $info->getExtension(), '|';
set_exception_handler(static function ($e) {});
set_exception_handler(static function ($e) {});
echo restore_exception_handler() ? 'restored' : 'failed';
