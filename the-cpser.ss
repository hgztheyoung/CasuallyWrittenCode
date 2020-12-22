(load "mk.scm")
(import (Framework match))
(import (Framework helpers))
(case-sensitive #t)
(print-gensym 'pretty)

(define-syntax λ
  (syntax-rules ()
    [(_ (x ...) body ...)
     (lambda (x ...) body ...)]))


(define (id x) x)

(define (k0 x) `(k ,x))

(define (cps p k)
  (match p
    [,sym (guard (symbol? sym)) (k sym)]
    [(λ (,x) ,body) (k `(λ (,x k) ,(cps body k0)))]
    [(,app ,rator)
     (cps app (λ (acode)
                (cps rator (λ (rcode)
                             `(,acode ,rcode
                                      ,(if (eq? k k0)
                                           'k
                                           (let ([sym (gensym)])
                                             `(λ (,sym)
                                                ,(k sym)))))))))]))


(define (atom? a)
  (not (or (null? a)
           (list? a))))

(define (subst p o n)
  (match p
    [() '()]
    [,a (guard (atom? a))
        (if (eq? a o) n a)]
    [(,[a] . ,[d]) `(,a . ,d)]))

(define (cps-i p k)
  (match p
    [,sym (guard (symbol? sym)) (k sym)]
    [(k ,sth) (k sth)]
    [(λ (,x k) ,body) (k `(λ (x) ,(cps-i body id)))]
    [(,app ,rator ,kont)
     (cps-i app (λ (acode)
                  (cps-i rator (λ (rcode)
                                 (let ([res `(,acode ,rcode)])
                                   (match kont
                                     [k (k res)]
                                     [(λ (,sym) ,body)
                                      (cps-i body (λ (bcode)
                                                  (k (subst bcode sym res))))]))))))]))

(define (cps p k)
  (match p
    [,sym (guard (symbol? sym)) (apply-k k sym)]
    [(λ (,x) ,body) (apply-k k `(λ (,x k) ,(cps body 'K-k0)))]
    [(,app ,rator)
     (cps app `(K-app-1 ,k ,rator))]))


(define (apply-k k code)
  (match k
    [K-id code]
    [K-k0 `(k ,code)]
    [(K-app-1 ,k ,rator)
     (cps rator `(K-app-0 ,k ,code))]
    [(K-app-0 ,k ,acode)
     `(,acode ,code
              ,(if (eq? k 'K-k0)
                   'k
                   (let ([sym (gensym)])
                     `(λ (,sym)
                        ,(apply-k k sym)))))]))



;;;;;
(define (cpsᵒ p k out)
  (conde
   [(symbolo p) (apply-kᵒ k p out)]
   [(fresh (x body bout)
      (== p `(λ (,x) ,body))
      (cpsᵒ body 'K-k0 bout)
      (apply-kᵒ k `(λ (,x k) ,bout) out))]
   [(fresh (app rator)
      (== p `(,app ,rator))
      (cpsᵒ app `(K-app-1 ,k ,rator) out))]))


(define (apply-kᵒ k-strcut code out)
  (conde
   [(== k-strcut 'K-id) (== code out)]
   [(== k-strcut 'K-k0) (== `(k ,code) out)]
   [(fresh (k rator)
      (== k-strcut `(K-app-1 ,k ,rator))
      (cpsᵒ rator `(K-app-0 ,k ,code) out))]
   [(fresh (k acode)
      (== k-strcut `(K-app-0 ,k ,acode))
      (conda
       [(== k 'K-k0) (== `(,acode ,code k) out)]
       [(fresh (sym appout)
          (symbolo sym)
          (apply-kᵒ k sym appout)
          (== `(,acode ,code (λ (,sym) ,appout)) out))]))]))


(define (verify-term program)
  (match program
    [,sym (guard (symbol? sym)) sym]
    [(λ (,x) ,[body]) (guard (symbol? x))
                      `(λ (,x) ,body)]
    [(,[app] ,[rator]) `(,app ,rator)]))


(define (verify-termᵒ out)
  (conde
   [(symbolo out)]
   [(fresh (x body)
      (== `(λ (,x) ,body) out)
      (symbolo x)
      (verify-termᵒ body))]
   [(fresh (app rator)
      (== `(,app ,rator) out)
      (verify-termᵒ app)
      (verify-termᵒ rator))]))


#!eof
(load "the-cpser.ss")


(cps `(p (λ (x) q)) 'K-id)

(cps-i `(p (λ (x k) (k q)) (λ (g10) g10))
 id)

(p (λ (x k) (k q)) (λ (#:g10) #:g10))


(cps `(p (λ (x) ((z p) q))) 'K-id)

(cps-i
 '(p (λ (x k) (z p (λ (g6) (g6 q k)))) (λ (g7) g7))
 id)


(cps '((λ (x) z) o) id)

(cps 'x id)

(run* (res)
  (cpsᵒ '(p (λ (x) ((z p) q)))  'K-id res ))

;;fail to terminate


(run* (res)
  (cpsᵒ res 'K-id '(z q (λ (x) x))))

(trace apply-kᵒ)
(trace cpsᵒ)
