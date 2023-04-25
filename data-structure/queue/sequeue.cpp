#include <iostream>

#define MaxSize 10

template <typename T>
class SeqQueue {
public:
    SeqQueue();
    ~SeqQueue();

public:
    bool EnQueue(const T& e);
    bool DeQueue(T &e);
    bool GetHead(T &e);
    void ClearQueue();
    
    void DispList();
    int ListLength();
    bool IsEmpty();
    bool IsFull();

private:
    T *m_data;
    int m_front;
    int m_rear;
};

template <typename T>
SeqQueue<T>::SeqQueue() {
    m_data = new T[MaxSize];
    m_front = 0;
    m_rear = 0;
}

template <typename T>
SeqQueue<T>::~SeqQueue() {
    delete[] m_data;
}

template <typename T>
bool SeqQueue<T>::EnQueue(const T &e) {
    if (IsFull() == true) {
        std::cout << "SeqQueue Full" << std::endl;
        return false;
    }
    m_data[m_rear] = e;
    m_rear++;
    return true;
}

template <typename T>
bool SeqQueue<T>::DeQueue(T &e) {
    if (IsEmpty() == true) {
        std::cout << "SeqQueue Empty" << std::endl;
        return false;
    }
    e = m_data[m_front];
    m_front++;
    return true;
}

template <typename T>
bool SeqQueue<T>::GetHead(T &e) {
    if (IsEmpty() == true) {
        std::cout << "SeqQueue Empty" << std::endl;
        return false;
    }
    e = m_data[m_front];
    return true;
}

template <class T>
void SeqQueue<T>::DispList() {
    for (int i = m_front; i < m_rear; i++) {
        std::cout << m_data[i] << " ";
    }
    std::cout << std::endl;
}

template <class T>
int SeqQueue<T>::ListLength() {
    return m_rear - m_front;
}

template <class T>
bool SeqQueue<T>::IsEmpty() {
    if (m_front == m_rear) {
        return true;
    }
    return false;
}

template <class T>
bool SeqQueue<T>::IsFull() {
    if (m_rear >= MaxSize) {
        return true;
    }
    return false;
}

template <class T>
void SeqQueue<T>::ClearQueue() {
    m_front = m_rear = 0;
}

int main(void) {
    SeqQueue<int> seqobj;
    seqobj.EnQueue(150);
    seqobj.EnQueue(200);
    seqobj.EnQueue(300);
    seqobj.EnQueue(400);
    seqobj.DispList();

    return 0;
}