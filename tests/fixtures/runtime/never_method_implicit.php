<?php
class Worker { public function stop(): never {} }
(new Worker)->stop();
