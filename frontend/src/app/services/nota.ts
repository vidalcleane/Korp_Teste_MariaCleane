import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';

export interface ItemNota {
  produto_id: number;
  quantidade: number;
}

export interface Nota {
  numero?: number;
  status: string;
  itens: ItemNota[];
}

@Injectable({
  providedIn: 'root'
})
export class NotaService {
  private readonly apiUrl = 'http://localhost:8082/notas';

  private notasSubject = new BehaviorSubject<Nota[]>([]);
  notas$: Observable<Nota[]> = this.notasSubject.asObservable();

  constructor(private http: HttpClient) {}

  carregarNotas(): void {
    this.http.get<Nota[]>(this.apiUrl).subscribe({
      next: (notas) => this.notasSubject.next(notas),
      error: (err) => console.error('Erro ao carregar notas', err)
    });
  }

  criarNota(itens: ItemNota[]): Observable<Nota> {
    return this.http.post<Nota>(this.apiUrl, { itens }).pipe(
      tap(() => this.carregarNotas())
    );
  }

  imprimirNota(numero: number): Observable<{ status: string }> {
    return this.http.post<{ status: string }>(`${this.apiUrl}/${numero}/imprimir`, {}).pipe(
      tap(() => this.carregarNotas())
    );
  }
}