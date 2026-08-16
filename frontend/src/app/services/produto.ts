import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';

export interface Produto {
  id?: number;
  codigo: string;
  descricao: string;
  saldo: number;
}

@Injectable({
  providedIn: 'root'
})
export class ProdutoService {
  private readonly apiUrl = 'http://localhost:8081/produtos';

  // BehaviorSubject guarda o estado atual da lista de produtos e avisa
  // automaticamente qualquer tela que estiver "assinando" sempre que
  // a lista mudar (ex: depois de cadastrar um produto novo).
  private produtosSubject = new BehaviorSubject<Produto[]>([]);
  produtos$: Observable<Produto[]> = this.produtosSubject.asObservable();

  constructor(private http: HttpClient) {}

  carregarProdutos(): void {
    this.http.get<Produto[]>(this.apiUrl).subscribe({
      next: (produtos) => this.produtosSubject.next(produtos),
      error: (err) => console.error('Erro ao carregar produtos', err)
    });
  }

  criarProduto(produto: Produto): Observable<Produto> {
    return this.http.post<Produto>(this.apiUrl, produto).pipe(
      tap(() => this.carregarProdutos())
    );
  }
}
