<?php
class Worker {
    public function finish(): void { echo 'done|'; }
    public static function stop(): never { throw new Exception('stop'); }
}
echo gettype((new Worker)->finish()), '|';
try { Worker::stop(); } catch (Exception $e) { echo $e->getMessage(), "\n"; }
