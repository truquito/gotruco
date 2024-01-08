# TODO

## BUGS

## LIMPIEZA DE CODIGO
  
- el eval de la flor esta dividido: el eval solo se llama cuando se juega flor comun
  mientras que los quiero/noquiero se llaman desde ahi mismo
  para el eval de contraflor/contrafloralresto no hay case (ya que se implementan
  en los quiero/noquiero)
- Hay codigo repetido entre noQuiero y mazo cuando niega la flor (codigo copiado)
- hay redundancia entre cantarFloresSiLasHay y cantarFlores
- hay redundancia entre p.Ronda.GetLaFlorMasAlta y 
    manojoConLaFlorGanadora, _, _ := p.Ronda.execCantarFlores(aPartirDe)

## SEGURIDAD

## PERFORMANCE

## DUDAS

