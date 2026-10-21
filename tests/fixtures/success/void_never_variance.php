<?php
interface Task { public function run(): void; }
abstract class BaseTask implements Task { abstract public function run(): void; }
trait StopTask { public function run(): never { throw new Exception('task'); } }
class RealTask extends BaseTask { use StopTask; }
try { (new RealTask)->run(); } catch (Exception $e) { echo $e->getMessage(), "\n"; }
