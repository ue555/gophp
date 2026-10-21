<?php
class PlainParent { public function run() { return 3; } }
class VoidChild extends PlainParent { public function run(): void {} }
class IntParent { public function run(): int { return 3; } }
class NeverChild extends IntParent { public function run(): never { throw new Exception('int'); } }
class MixedParent { public function run(): mixed { return 3; } }
class NeverMixedChild extends MixedParent { public function run(): never { throw new Exception('mixed'); } }
echo gettype((new VoidChild)->run()), '|';
try { (new NeverChild)->run(); } catch (Exception $e) { echo $e->getMessage(), '|'; }
try { (new NeverMixedChild)->run(); } catch (Exception $e) { echo $e->getMessage(), "\n"; }
