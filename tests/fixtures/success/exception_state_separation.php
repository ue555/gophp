<?php
try { throw new Exception('one'); } catch (Exception $first) {}
try { throw new Error('two'); } catch (Error $second) {
    echo ($second->getPrevious() === null ? 'clear' : 'leaked'), '|';
    echo $second->getMessage(), "\n";
}
