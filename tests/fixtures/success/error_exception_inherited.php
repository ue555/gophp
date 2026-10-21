<?php
class ViewProblem extends ErrorException {}
$previous=new Exception('cause');
$e=new ViewProblem('view',7,E_WARNING,'template.php',42,$previous);
echo $e->getMessage(),':',$e->getCode(),':',$e->getSeverity(),':',$e->getFile(),':',$e->getLine(),':',$e->getPrevious()->getMessage(),"\n";
