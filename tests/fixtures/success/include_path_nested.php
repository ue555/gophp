<?php
set_include_path(__DIR__ . '/../support/include_path/first');
function loadNeighbor() { return include __DIR__ . '/../support/include_path/second/nested.inc'; }
class Loader { public function load() { return require __DIR__ . '/../support/include_path/second/nested.inc'; } }
echo loadNeighbor(), '|', (new Loader)->load(), "\n";
