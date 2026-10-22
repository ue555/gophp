<?php
namespace WorkingDirectory;
$cwd = getcwd();
echo gettype($cwd), '|', is_dir($cwd) ? 'directory' : 'missing', '|', $cwd === \getcwd() ? 'stable' : 'changed', "\n";
foreach ([[1], ['unknown' => 1]] as $args) {
    try { getcwd(...$args); }
    catch (\Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
}
echo getcwd() === $cwd ? "unchanged\n" : "changed\n";
