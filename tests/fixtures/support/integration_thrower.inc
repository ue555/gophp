<?php
function integratedThrower(): Closure {
    return function (): void { throw new Exception('inside'); };
}
return 'ready';
